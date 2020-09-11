from django.contrib.auth import views as auth_views
from django.urls import reverse


class LoginView(auth_views.LoginView):
    template_name = 'accounts/login.html'


class LogoutView(auth_views.LogoutView):
    template_name = 'accounts/logged_out.html'
