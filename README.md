How to run this program
===============

### with docker
0. Clone this repository
1. ```docker build -t goplayground .```
2. ```docker run -p 2345:8000 -d goplayground```
3. open localhost:2345 on your browser

### with goland ide
1. Select get from VCS
2. Insert this repo url
3. comment ```	var addr = "0.0.0.0:8000"``` and uncomment ```var addr = "localhost:8000"```
4. Run the project
