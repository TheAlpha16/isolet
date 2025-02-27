from database import Database
from dotenv import load_dotenv
import json

load_dotenv()
DB = Database()

class Hint:
    def __init__(self, hint, cost, visible=False):
        self.hint = hint
        self.cost = cost
        self.visible = visible

        self.__linter()

    def __linter(self):
        fields = [
            {"object": self.hint, "type": str, "optional": False},
            {"object": self.cost, "type": int, "optional": False},
            {"object": self.visible, "type": bool, "optional": True}
        ]
            
        for field in fields:
            if not isinstance(field["object"], field["type"]):
                raise ValueError(f"Invalid type for {field['object']} in {self.__class__.__name__} for {field['object']}")

            if not field["optional"]:
                if not field["object"]:
                    raise ValueError(f"Missing {field['object']} in {self.__class__.__name__} for {field['object']}")

        if self.cost < 0:
            raise ValueError(f"Cost must be greater than 0 for {field['object']}")
        
    def __repr__(self):
        return f"Hint({self.hint}, {self.cost}, {self.visible})"
    
    def hint(self):
        return self.hint
    
    def cost(self):
        return self.cost
    
    def visible(self):
        return self.visible

class Challenge:
    def __init__(self, chall_name, category, type, points, prompt=None, flag=None, files=[], requirements=[], hints=[], author=None, tags=[], links=[], image=None, registry=None, deployment=None, port=None, subd=None, cpu=None, mem=None):
        self.chall_name = chall_name
        self.category = category
        self.prompt = prompt
        self.flag = flag
        self.type = type
        self.points = points
        self.files = files
        self.requirements = requirements
        self.hint_objects = hints
        self.author = author
        self.tags = tags
        self.links = links
        self.image_metadata = {
            "image": image,
            "registry": registry,
            "deployment": deployment,
            "port": port,
            "subd": subd,
            "cpu": cpu,
            "mem": mem
        }

        self.fields = [
            {"object": self.chall_name, "type": str, "optional": False, "key": "chall_name"},
            {"object": self.category, "type": str, "optional": False, "key": "category"},
            {"object": self.prompt, "type": str, "optional": True, "default": "", "key": "prompt"},
            {"object": self.flag, "type": str, "optional": True, "default": "", "key": "flag"},
            {"object": self.type, "type": str, "optional": False, "key": "type"},
            {"object": self.points, "type": int, "optional": False, "key": "points"},
            {"object": self.files, "type": list, "optional": True, "default": [], "key": "files"},
            {"object": self.requirements, "type": list, "optional": True, "default": [], "key": "requirements"},
            {"object": self.hint_objects, "type": list, "optional": True, "default": [], "key": "hint_objects"},
            {"object": self.author, "type": str, "optional": True, "default": "anonymous", "key": "author"},
            {"object": self.tags, "type": list, "optional": True, "default": [], "key": "tags"},
            {"object": self.links, "type": list, "optional": True, "default": [], "key": "links"},
            {"object": self.image_metadata, "type": dict, "optional": True, "default": {}, "key": "image_metadata"},
        ]

        self.__linter()

    def __linter(self):
        for field in self.fields:
            if not field["optional"]:
                if not field["object"]:
                    raise ValueError(f"Missing {field['key']}")
                if not isinstance(field["object"], field["type"]):
                    raise ValueError(f"Invalid type for {field['key']}")
            else:
                if field["object"]:
                    if not isinstance(field["object"], field["type"]):
                        raise ValueError(f"Invalid type for {field['key']}")

        if not self.type in ['static', 'dynamic', 'on-demand']:
            raise ValueError("Invalid challenge type")
        
        if self.type != 'static':
            if not self.image_metadata["image"]:
                raise ValueError(f"Missing image name for {self.chall_name} (dynamic or on-demand)")
            if not self.image_metadata["port"]:
                raise ValueError(f"Missing port metadata for {self.chall_name} (dynamic or on-demand)")
            if not self.image_metadata["deployment"]:
                raise ValueError(f"Missing deployment metadata for {self.chall_name} (dynamic or on-demand)")
            elif not self.image_metadata["deployment"] in ['http', 'ssh', 'nc']:
                raise ValueError(f"Invalid deployment method for {self.chall_name}, should be in ('http', 'ssh', 'nc')")

        if self.points < 0:
            raise ValueError("Points must be greater or equal than 0")
        
        self.hint_objects = [Hint(hint["hint"], hint["cost"], hint.get("visible", False)) for hint in self.hint_objects]

    def props(self):
        return {field["key"]: field["object"] for field in self.fields}

    def create(self):
        DB.new_challenge(self.props())

def main():
    challenges = json.loads(open('challs.json', 'r').read())
    final_challenges = []
    got_till_now = []
    processed = True

    while len(final_challenges) < len(challenges):
        if not processed:
            raise ValueError("Circular dependency detected, check requirements")
        
        processed = False
        for challenge in challenges:
            if challenge["chall_name"] in got_till_now:
                continue

            if not challenge["requirements"]:
                processed = True
                final_challenges.append(challenge)
                got_till_now.append(challenge["chall_name"])
            else:
                if all([req in got_till_now for req in challenge["requirements"]]):
                    processed = True
                    final_challenges.append(challenge)
                    got_till_now.append(challenge["chall_name"])

    for challenge in final_challenges:
        Challenge(**challenge).create()

if __name__ == "__main__":
    main()