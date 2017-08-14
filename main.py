from flask import Flask, jsonify, render_template, Response, request
from flask_mysqldb import MySQL
import json
import os.path
import uuid

class Event:
    name = ""
    attributes = {}
    
    def __init__(self, json):
        self.name = json["name"]
        self.attributes = json["attributes"]
    
    def attributes_json(self):
        return json.dumps(self.attributes, separators=(",", ":"))

api_key = ""
app = Flask(__name__)
sql = MySQL()

app.config["MYSQL_USER"] = "root"
app.config["MYSQL_PASSWORD"] = "password"
app.config["MYSQL_DB"] = "fuji"
app.config["MYSQL_HOST"] = "localhost"
sql.init_app(app)

@app.route("/dashboard")
def dashboard():
    return render_template("dashboard.html")

@app.route("/data")
def data():
    pages = {}
    
    cursor = sql.connection.cursor()
    cursor.execute("SELECT * FROM events WHERE name = 'Content View'")
    result = cursor.fetchall()
    
    for row in result:
        page = json.loads(row[2])["page"]
        
        if page in pages:
            pages[page] += 1
        else:
            pages[page] = 1
    
    contents = "page,count\n"
    
    for page in pages:
        contents += "%s,%d\n" % (page, pages[page])
    
    return Response(contents, mimetype="text/csv")

@app.route("/events", methods = ["POST"])
def events():
    if request.headers.get("Authorization") != api_key:
        return jsonify({"message": "Unauthorized"}), 401
    
    event = Event(request.json)
    query = "INSERT INTO events (name, attributes) VALUES (%s, %s)"
    
    cursor = sql.connection.cursor()
    cursor.execute(query, (event.name, event.attributes_json()))
    sql.connection.commit()
    
    return jsonify({"message": "Hello"})

if __name__ == "__main__":
    if os.path.isfile("config.json"):
        with open("config.json", "r") as f:
            data = json.load(f)
            api_key = data["API_KEY"]
    else:
        with open("config.json", "w") as f:
            api_key = str(uuid.uuid4()).replace("-", "")
            data = {"API_KEY": api_key}
            json.dump(data, f)
    
    print("API key: %s" % api_key)
    app.run()
