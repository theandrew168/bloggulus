from pyramid.response import Response


def hello_world(request):
    return Response('<body><h1>Hello World!</h1></body>')
