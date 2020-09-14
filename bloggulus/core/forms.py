from django import forms
from django.core.exceptions import ValidationError
from django.utils.translation import gettext_lazy as _
import feedparser

from .models import Feed, Post


class RSSFeedForm(forms.Form):
    url = forms.URLField(label='Feed URL')

    def clean(self):
        super().clean()

        url = self.cleaned_data.get('url')
        if not url:
            # short-circuit if URL isn't valid
            return self.cleaned_data

        d = feedparser.parse(url)
        if 'feed' not in d:
            raise ValidationError(_('Invalid RSS feed'), code='invalid')
        if 'entries' not in d:
            raise ValidationError(_('Invalid RSS feed'), code='invalid')

        feed = d['feed']
        if 'title' not in feed:
            raise ValidationError(_('Invalid RSS feed'), code='invalid')
        if 'updated' not in feed:
            raise ValidationError(_('Invalid RSS feed'), code='invalid')

        return self.cleaned_data

    def add_feed(self):
        url = self.cleaned_data['url']

        d = feedparser.parse(url)
        feed = d['feed']
        posts = d['entries']

        title = feed['title']
        updated = feed['updated']

        existing = Feed.objects.filter(url=url)
        if existing.exists():
            return
 
        f = Feed(title=title, url=url, updated=updated)
        f.save()

        for post in posts:
            title = post.get('title')
            url = post.get('link')
            updated = post.get('updated')

            if not all([title, url, updated]):
                continue
 
            p = Post(feed=f, title=title, url=url, updated=updated)
            p.save()
