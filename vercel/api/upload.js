const express = require('express');
const fileUpload = require('express-fileupload');
const fs = require('fs');
const os = require('os');
const path = require('path');

const app = express();

app.use(fileUpload());

app.get('/api/hello', (req, res) => {
    res.send('Hello World!');
});

app.post('/api/upload', (req, res) => {
    if (!req.files || Object.keys(req.files).length === 0) {
        return res.status(400).send('No files were uploaded.');
    }

    // The name of the input field is used to retrieve the uploaded file
    let uploadedFile = req.files.file;

    // Get the current date in ISO format
    let folderPath = new Date().toISOString().split('T')[0];
    let userHome = os.homedir();
    let screenshotsFolder = path.join(userHome, 'screenshots');
    let folderFullPath = path.join(screenshotsFolder, folderPath);

    // Create the screenshots folder if it doesn't exist
    if (!fs.existsSync(screenshotsFolder)) {
        fs.mkdirSync(screenshotsFolder, { recursive: true });
    }

    // Create the folder with the current date inside the screenshots folder
    if (!fs.existsSync(folderFullPath)) {
        fs.mkdirSync(folderFullPath, { recursive: true });
    }

    let filePath = path.join(folderFullPath, uploadedFile.name);

    // Use the mv() method to place the file on the server
    uploadedFile.mv(filePath, (err) => {
        if (err) {
            console.log('Error saving file:', err);
            return res.status(500).send('Error saving file.');
        }

        console.log('File saved:', filePath);
        res.send('File uploaded and saved.');
    });
});

module.exports = app;
