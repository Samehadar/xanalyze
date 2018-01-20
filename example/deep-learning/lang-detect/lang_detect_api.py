from __future__ import unicode_literals
from flask import Flask, jsonify, request, send_from_directory
import logging
from logging import FileHandler
import time
import data_helper
import keras
import numpy as np
import defs

t = time.time()
if __name__ == '__main__':
  print("Time is now {}".format(t))

app = Flask(__name__)
application = app

try:
  file_handler = FileHandler("lang_detection.log", "a")
  file_handler.setLevel(logging.WARNING)
  app.logger.addHandler(file_handler)
  model = keras.models.load_model('./save_tmp.h5')
  if __name__ == '__main__':
    print("Took {}".format(time.time() - t))
except:
  if __name__ == '__main__':
    print("could not start file logging")

@app.route('/', methods=['GET'])
def browse_default():
  try:
    return send_from_directory('ui', 'index.html')
  except Exception as e:
    return e.message

@app.route('/<path:path>', methods=['GET'])
def staticx(path):
   return send_from_directory('ui', path)

@app.route('/api/v1/lang/detection', methods=['GET'])
def detect_url():

  f_url = request.args.get('file')
  x = np.array([data_helper.turn_url_to_vector(f_url, file_vector_size=2048)])
  y = model.predict(x)
  result = model.predict_proba(x)
  return jsonify({"result": result[0].tolist()})

@app.route('/api/v1/lang/detection', methods=['POST'])
def detect_text():

  text = request.get_json()['content']
  x = np.array([data_helper.turn_text_to_vector(text, file_vector_size=2048)])
  y = model.predict(x)
  result = [100*a for a in model.predict_proba(x)] # tunr to percentage
  return jsonify({"result": result[0].tolist()})

if __name__ == '__main__':
  app.run(debug = True, use_reloader=False)

