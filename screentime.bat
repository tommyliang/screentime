@echo off

:loop
rem Get the current time with colons replaced by dashes
set "NOW=%TIME:~0,2%-%TIME:~3,2%-%TIME:~6,2%"
set "NOW=%NOW::=-%"

rem Get today's date and format it as YYYY-MM-DD
for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set "datetime=%%I"
set "TODAY=%datetime:~0,4%-%datetime:~4,2%-%datetime:~6,2%"

rem Create a subfolder with today's date if it does not exist
if not exist "%TODAY%" mkdir "%TODAY%"

rem Capture a screenshot and save it to the date subfolder with the current time in the filename
set SAVE_TO_FILE=%TODAY%/screenshot_%NOW%.jpg
ffmpeg -f gdigrab -framerate 1 -offset_x 0 -offset_y 0 -video_size 1920x1080 -i desktop -frames:v 1 -q:v 31 -loglevel quiet "%SAVE_TO_FILE%"

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
