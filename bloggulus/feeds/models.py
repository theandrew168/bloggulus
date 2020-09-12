from django.contrib.auth.models import User
from django.db import models


class Feed(models.Model):
    users = models.ManyToManyField(User, related_name='feeds')
    title = models.CharField(max_length=200)
    url = models.URLField(unique=True)
    updated = models.DateTimeField()

    def __str__(self):
        return self.title


class Post(models.Model):
    feed = models.ForeignKey(Feed, on_delete=models.CASCADE)
    title = models.CharField(max_length=200)
    url = models.URLField(unique=True)
    updated = models.DateTimeField()

    def __str__(self):
        return self.title
