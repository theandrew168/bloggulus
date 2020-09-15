from datetime import datetime
from time import mktime

from django.contrib.auth import authenticate, login
from django.contrib.auth.decorators import login_required
from django.contrib.auth.forms import UserCreationForm
from django.http import HttpResponseRedirect
from django.template.response import TemplateResponse
from django.urls import reverse
from django.views.decorators.http import require_http_methods
import feedparser
import pytz

from .forms import RSSFeedForm
from .models import Feed, Post


@require_http_methods(['GET', 'HEAD'])
def index(request):
    posts = Post.objects.order_by('-updated')[:20]
    return TemplateResponse(request, 'core/posts.html', {'posts': posts})


@login_required
@require_http_methods(['GET', 'HEAD'])
def posts(request):
    feeds = Feed.objects.filter(users=request.user)
    posts = Post.objects.filter(feed__in=feeds).order_by('-updated')[:20]
    return TemplateResponse(request, 'core/posts.html', {'posts': posts})


@require_http_methods(['GET', 'HEAD', 'POST'])
def register(request):
    if request.method == 'POST':
        form = UserCreationForm(request.POST)
        if form.is_valid():
            # create the user
            form.save()

            # log the user in
            username = form.cleaned_data['username']
            password = form.cleaned_data['password1']
            user = authenticate(username=username, password=password)
            login(request, user)

            return HttpResponseRedirect(reverse('core:posts'))
    else:
        form = UserCreationForm()

    return TemplateResponse(request, 'core/register.html', {'form': form})


@login_required
@require_http_methods(['GET', 'HEAD', 'POST'])
def profile(request):
    if request.method == 'POST':
        form = RSSFeedForm(request.POST)
        if form.is_valid():
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
                f[0].users.add(request.user)
                return HttpResponseRedirect(reverse('core:posts'))
 
            f = Feed(title=title, url=url, updated=updated)
            f.save()
            f.users.add(request.user)

            for post in posts:
                title = post.get('title')
                url = post.get('link')
                updated = datetime.fromtimestamp(mktime(post.get('updated_parsed')))
                updated = pytz.UTC.localize(updated)

                # skip posts with incomplete data
                if None in [title, url, updated]:
                    continue
 
                p = Post(feed=f, title=title, url=url, updated=updated)
                p.save()

        return HttpResponseRedirect(reverse('core:posts'))
    else:
        form = RSSFeedForm()

    return TemplateResponse(request, 'core/profile.html', {'form': form})
            

@require_http_methods(['GET', 'HEAD'])
def about(request):
    return TemplateResponse(request, 'core/about.html')


@require_http_methods(['GET', 'HEAD'])
def terms(request):
    return TemplateResponse(request, 'core/terms.html')


@require_http_methods(['GET', 'HEAD'])
def privacy(request):
    return TemplateResponse(request, 'core/privacy.html')
