from datetime import datetime
from time import mktime

from django.core.management.base import BaseCommand, CommandError
import feedparser
import pytz

from bloggulus.core.models import Feed, Post


class Command(BaseCommand):
    help = 'Check RSS feeds for new posts'

    def handle(self, *args, **options):
        for f in Feed.objects.all():
            self.stdout.write('Syncing feed: {}'.format(f.title))
            d = feedparser.parse(f.url)

            feed = d['feed']
            updated = datetime.fromtimestamp(mktime(feed['updated_parsed']))
            updated = pytz.UTC.localize(updated)

            # check when feed was last updated
#            if updated <= f.updated:
#                self.stdout.write('  not updated recently')
#                continue

            # iterate over each post in the feed
            posts = d['entries']
            for post in posts:
                title = post.get('title')
                url = post.get('link')
                updated = datetime.fromtimestamp(mktime(post.get('updated_parsed')))
                updated = pytz.UTC.localize(updated)
    
                # skip if required info isn't present
                if None in [title, url, updated]:
                    continue

                # check if post exists (by URL) and create if not
                try:
                    p = Post.objects.get(url=url)
                except Post.DoesNotExist:
                    self.stdout.write('  Adding post: {}'.format(title))
                    p = Post(feed=f, title=title, url=url, updated=updated)
                    p.save()
                else:
                    self.stdout.write('  Post exists: {}'.format(p.title))

            # update feed's updated datetime
            f.updated = updated
            f.save()

        self.stdout.write(self.style.SUCCESS('Successfully synced feeds!'))
