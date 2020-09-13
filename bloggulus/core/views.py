from django.contrib.auth import views as auth_views
from django.views import generic

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


class ProfileView(generic.TemplateView):
    template_name = 'core/profile.html'


class AboutView(generic.TemplateView):
    template_name = 'core/about.html'


class TermsView(generic.TemplateView):
    template_name = 'core/terms.html'


class PrivacyView(generic.TemplateView):
    template_name = 'core/privacy.html'


# def process(request):
#     url = request.POST['feed'].strip()
# 
#     d = feedparser.parse(url)
#     feed = d['feed']
#     posts = d['entries']
# 
#     title = feed['title']
#     updated = feed['updated']
# 
#     try:
#         existing = Feed.objects.get(url=url)
#     except Feed.DoesNotExist:
#         pass
#     else:
#         print('skipping existing feed!')
#         return HttpResponseRedirect(reverse('feeds:index'))
# 
#     f = Feed(title=title, url=url, updated=updated)
#     f.save()
# 
#     for post in posts:
#         title = post['title']
#         url = post['link']
#         updated = post['updated']
#         p = Post(feed=f, title=title, url=url, updated=updated)
#         p.save()
# 
#     return HttpResponseRedirect(reverse('feeds:index'))
