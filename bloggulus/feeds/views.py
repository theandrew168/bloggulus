from django.http import HttpResponseRedirect
from django.urls import reverse
from django.views import generic
import feedparser

from .models import Feed, Post


class IndexView(generic.ListView):
    template_name = 'feeds/index.html'
    context_object_name = 'latest_posts'

    def get_queryset(self):
        return Post.objects.order_by('-updated')[:10]


class AddView(generic.TemplateView):
    template_name = 'feeds/add.html'


def process(request):
    url = request.POST['feed_url']

    d = feedparser.parse(url)
    feed = d['feed']
    posts = d['entries']

    title = feed['title']
    updated = feed['updated']

    try:
        existing = Feed.objects.get(url=url)
    except Feed.DoesNotExist:
        pass
    else:
        print('skipping existing feed!')
        return HttpResponseRedirect(reverse('feeds:index'))

    f = Feed(title=title, url=url, updated=updated)
    f.save()

    for post in posts:
        title = post['title']
        url = post['link']
        updated = post['updated']
        p = Post(feed=f, title=title, url=url, updated=updated)
        p.save()

    return HttpResponseRedirect(reverse('feeds:index'))
