import sys
import requests
import time
import os

def upload_file(url, file_path):
    try:
        file = open(file_path, 'rb')
        response = requests.post(url, files={'file': file})
        file.close()  # Close the file explicitly
        if response.status_code == 200:
            print("File {} uploaded successfully!".format(file_path))

            # Wait for the file to be released before removing it
            while True:
                try:
                    os.remove(file_path)
                    print('File removed:', file_path)
                    break  # Exit the loop if the file is successfully removed
                except PermissionError:
                    print('File is still in use, waiting...')
                    time.sleep(1)  # Wait for 1 second before retrying
        else:
            print(f"File upload failed with status code: {response.status_code}")
    except IOError as e:
        print(f"Error opening file: {e}")
    except requests.exceptions.RequestException as e:
        print(f"Error uploading file: {e}")

if __name__ == '__main__':
    if len(sys.argv) < 3:
        print("Usage: python upload_script.py <url> <file_path>")
    else:
        url = sys.argv[1]
        file_path = sys.argv[2]
        upload_file(url, file_path)
