@echo off

set /p PROJ_NAME=Input a project name:

xcopy %cd%\skeleton\*.* %GOPATH%\src\%PROJ_NAME%\ /s /e /i /y

cd %GOPATH%\src\%PROJ_NAME%\conf

set temp=framework.tmp

for /r %%k in (framework.conf.dev.xml framework.conf.test.xml framework.conf.product.xml) do (
	
	if exist "%%k" (
		cd.>"%temp%"

		for /f "tokens=1* delims=:" %%i in ('findstr /i /n .* "%%k"') do (
			set str=%%j
			setlocal enabledelayedexpansion
			if not "!str!"=="" set "str=!str:beekeeper=%PROJ_NAME%!"
			>>"%temp%" echo.!str!
			endlocal
		)
		
		move "%temp%" "%%k"
	)
)

cd %GOPATH%\src\%PROJ_NAME%

set temp=main.tmp

for /r %%k in (main.go) do (
	
	if exist "%%k" (
		cd.>"%temp%"

		for /f "tokens=1* delims=:" %%i in ('findstr /i /n .* "%%k"') do (
			set str=%%j
			setlocal enabledelayedexpansion
			if not "!str!"=="" set "str=!str:github.com/boringding/beekeeper/skeleton/handlers=%PROJ_NAME%/handlers!"
			>>"%temp%" echo.!str!
			endlocal
		)
		
		move "%temp%" "%%k"
	)
)
