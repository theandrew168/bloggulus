from datetime import datetime
from time import mktime

from django.contrib.auth import authenticate, login
from django.contrib.auth import views as auth_views
from django.contrib.auth.forms import UserCreationForm
from django.contrib.auth.mixins import LoginRequiredMixin
from django.urls import reverse_lazy
from django.views import generic
import feedparser
import pytz

from .forms import RSSFeedForm
from .models import Feed, Post


class IndexView(generic.ListView):
    template_name = 'core/posts.html'
    context_object_name = 'posts'

    # display recent posts from ALL feeds
    def get_queryset(self):
        return Post.objects.order_by('-updated')[:20]


class PostsView(LoginRequiredMixin, generic.ListView):
    template_name = 'core/posts.html'
    context_object_name = 'posts'

    # display recent posts from the current user's feeds
    def get_queryset(self):
        feeds = Feed.objects.filter(users=self.request.user)
        return Post.objects.filter(feed__in=feeds).order_by('-updated')[:20]


class LoginView(auth_views.LoginView):
    template_name = 'core/login.html'


class LogoutView(auth_views.LogoutView):
    pass


class RegisterView(generic.FormView):
    template_name = 'core/register.html'
    form_class = UserCreationForm
    success_url = reverse_lazy('core:posts')

    def form_valid(self, form):
        form.save()
        username = form.cleaned_data['username']
        password = form.cleaned_data['password1']
        user = authenticate(username=username, password=password)
        login(self.request, user)
        return super().form_valid(form)


# TODO: would a CreateView or UpdateView be better here?
class ProfileView(LoginRequiredMixin, generic.FormView):
    template_name = 'core/profile.html'
    form_class = RSSFeedForm
    success_url = reverse_lazy('core:posts')

    def form_valid(self, form):
        url = form.cleaned_data['url']

        d = feedparser.parse(url)
        feed = d['feed']
        posts = d['entries']

        title = feed['title']
        updated = datetime.fromtimestamp(mktime(feed['updated_parsed']))
        updated = pytz.UTC.localize(updated)

        # TODO: better to user .get() and catch DNE here?
        f = Feed.objects.filter(url=url)
        if f.exists():
            f[0].users.add(self.request.user)
            return super().form_valid(form)
 
        f = Feed(title=title, url=url, updated=updated)
        f.save()
        f.users.add(self.request.user)

        for post in posts:
            title = post.get('title')
            url = post.get('link')
            updated = datetime.fromtimestamp(mktime(post.get('updated_parsed')))
            updated = pytz.UTC.localize(updated)

            if not all([title, url, updated]):
                continue
 
            p = Post(feed=f, title=title, url=url, updated=updated)
            p.save()

        return super().form_valid(form)


class AboutView(generic.TemplateView):
    template_name = 'core/about.html'


class TermsView(generic.TemplateView):
    template_name = 'core/terms.html'


class PrivacyView(generic.TemplateView):
    template_name = 'core/privacy.html'
