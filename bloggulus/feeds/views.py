from django.views import generic

from .models import Entry


class IndexView(generic.ListView):
    template_name = 'feeds/index.html'
    context_object_name = 'latest_entries'

    def get_queryset(self):
        return Entry.objects.order_by('-updated')[:10]
