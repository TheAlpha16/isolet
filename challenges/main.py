from .database import Database

# {
#     "chall_name": "Ergonomic clear-thinking neural-net",
#     "category": "web",
#     "prompt": "Raise pretty walk impact.",
#     "flag": "w6H$zc%o*p%W",
#     "type": "on-demand",
#     "points": 179,
#     "files": [
#         "sure.txt",
#         "wall.txt",
#         "section.txt"
#     ],
#     "hints": [
#         "Same talk reach ability prove just.",
#         "Identify study again economic.",
#         "Actually view four."
#     ],
#     "author": "jphillips",
#     "tags": [
#         "morning",
#         "form",
#         "power"
#     ],
#     "port": 49912,
#     "subd": "srv-77.wilson-lynch.com",
#     "cpu": 10,
#     "mem": 3
# }

class Challenge:
    def __init__(self, chall_name, category, chall_type, points, prompt=None, flag=None, files=[], hints=[], author=None, tags=[], port=None, subd=None, cpu=None, mem=None):
        self.chall_name = chall_name
        self.category = category
        self.prompt = prompt
        self.flag = flag
        self.chall_type = chall_type
        self.points = points
        self.files = files
        self.hints = hints
        self.author = author
        self.tags = tags
        self.port = port
        self.subd = subd
        self.cpu = cpu
        self.mem = mem

    def __lint(self):
        if self.chall_name is None:
            return "chall_name is required"
        if self.category is None:
            return "category is required"
        if self.chall_type is None:
            return "chall_type is required"
        if self.points is None:
            return "points is required"
        return None