from django.views import generic

from .models import Post


class IndexView(generic.ListView):
    template_name = 'feeds/index.html'
    context_object_name = 'latest_posts'

    def get_queryset(self):
        return Post.objects.order_by('-updated')[:10]
