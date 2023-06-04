const { google } = require('googleapis');
const fs = require('fs');
const path = require('path');

const key = {
    client_email: "245140179831-fj7m2mtmdskm3suvibcdq5orp0jei7kv.apps.googleusercontent.com",
    private_key: "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"
};

const jwtClient = new google.auth.JWT(
  key.client_email,
  null,
  key.private_key,
  ['https://www.googleapis.com/auth/drive'],
  null
);

jwtClient.authorize(function(err, tokens) {
  if (err) {
    console.log(err);
    return;
  } 

  const drive = google.drive({ version: 'v3', auth: jwtClient });

  const fileMetadata = {
    'name': 'your_file_name.jpg'
  };
  
  const media = {
    mimeType: 'image/jpeg',
    body: fs.createReadStream(path.join(__dirname, 'screenshot_22-17-58.jpg'))
  };

  drive.files.create({
    resource: fileMetadata,
    media: media,
    fields: 'id'
  }, (err, file) => {
    if (err) {
      console.error(err);
    } else {
      console.log('File Id: ', file.id);
    }
  });
});
