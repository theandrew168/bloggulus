from django.contrib.auth import views as auth_views
from django.contrib.auth.mixins import LoginRequiredMixin
from django.urls import reverse_lazy
from django.views import generic

from .forms import RSSFeedForm
from .models import Feed, Post


class IndexView(generic.ListView):
    template_name = 'core/index.html'
    context_object_name = 'posts'

    # TODO: different posts if user is logged in
    # grab only posts from the feeds they follow
    def get_queryset(self):
        return Post.objects.order_by('-updated')[:10]


class LoginView(auth_views.LoginView):
    template_name = 'core/login.html'


class LogoutView(auth_views.LogoutView):
    pass


class ProfileView(LoginRequiredMixin, generic.FormView):
    template_name = 'core/profile.html'
    form_class = RSSFeedForm
    success_url = reverse_lazy('core:index')

    def form_valid(self, form):
        form.add_feed()
        return super().form_valid(form)


class AboutView(generic.TemplateView):
    template_name = 'core/about.html'


class TermsView(generic.TemplateView):
    template_name = 'core/terms.html'


class PrivacyView(generic.TemplateView):
    template_name = 'core/privacy.html'
