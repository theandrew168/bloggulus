from django.conf import settings
from django.contrib import admin
from django.urls import include, path

urlpatterns = [
    path('', include('bloggulus.core.urls')),
    path('feeds/', include('bloggulus.feeds.urls')),
    path('accounts/', include('bloggulus.accounts.urls')),
    path('admin/', admin.site.urls),
]

if settings.DEBUG:
    import debug_toolbar
    urlpatterns.insert(0, path('__debug__/', include(debug_toolbar.urls)))
