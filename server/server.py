from flask import Flask, request
import os
from datetime import date

app = Flask(__name__)

@app.route('/upload', methods=['POST'])
def upload():
    if 'file' not in request.files:
        return 'No file uploaded.', 400

    file = request.files['file']
    if file.filename == '':
        return 'No file selected.', 400

    filename = file.filename
    folder_path = date.today().isoformat()
    user_home = os.path.expanduser('~')
    screenshots_folder = os.path.join(user_home, 'screenshots')
    folder_full_path = os.path.join(screenshots_folder, folder_path)

    # Create the screenshots folder if it doesn't exist
    if not os.path.exists(screenshots_folder):
        os.makedirs(screenshots_folder)

    # Create the folder with the current date inside the screenshots folder
    if not os.path.exists(folder_full_path):
        os.makedirs(folder_full_path)

    file_path = os.path.join(folder_full_path, filename)

    try:
        file.save(file_path)
        print('File saved:', file_path)
        return 'File uploaded and saved.'
    except Exception as e:
        print('Error saving file:', str(e))
        return 'Error saving file.', 500

if __name__ == '__main__':
    app.run(host='0.0.0.0')
