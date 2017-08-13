from flask import Flask, jsonify, request
from flask_mysqldb import MySQL
import json

class Event:
    name = ""
    attributes = {}
    
    def __init__(self, json):
        self.name = json["name"]
        self.attributes = json["attributes"]
    
    def attributes_json(self):
        return json.dumps(self.attributes, separators=(",", ":"))

app = Flask(__name__)
sql = MySQL()

app.config["MYSQL_USER"] = "root"
app.config["MYSQL_PASSWORD"] = "password"
app.config["MYSQL_DB"] = "fuji"
app.config["MYSQL_HOST"] = "localhost"
sql.init_app(app)

@app.route("/events", methods = ["POST"])
def events():
    event = Event(request.json)
    query = "INSERT INTO events (name, attributes) VALUES (%s, %s)"
    
    cursor = sql.connection.cursor()
    cursor.execute(query, (event.name, event.attributes_json()))
    sql.connection.commit()
    
    return jsonify({"message": "Hello"})

if __name__ == "__main__":
    app.run(debug=True)
