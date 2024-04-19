# Cars Viewer

## 1. Rut the Cars API

Cars API zip file can be found in this repository.  

```
./api.zip
```

Unpack it and install required packages following the instructions given in Cars API's README.  
Run the Cars API **using a separate terminal**. Instructions can be found in the Cars API README.  
Note that Cars API must be run in it's default port 3000 for this Cars Viewer to work.

## 2. Run the Cars Viewer server

Run Cars Viewer server with command `go run .`. Cars Viewer is run in port 8080.  
Server up when this message is shown in the terminal:

```
$ go run .  
[Date and time] Staring server on port 8080...
```

Server is shut down by pressing `ctrl` + `C` in the terminal.

## 3. Interface usage

Make sure both servers are up and running.  
Open a browser and type  `http://localhost:8080` to the address bar and press enter.

The interface is now open.

### Home page

In the home page there are search filtering options for cars manufacturer and category.  
One option may be chosen in each filter. One of the filters may also be left in it's default value (empty).  
Search is done by pressing `Search` button.

There is a recommendation banner in the bottom of the home page that shows a recommendation of a car that might interest the user.
Recommendations are based on cookie information collected when user does searches and checks the specifications of cars.  
First recommendation is given at random.
By pressing `See specifications` button user can see the specifications of the recommended car.

### Search results page

The search results page shows the results of the search. Results are given in order of fitness to filters used: At the top are shown results that fit both selected filters and above are the search results that fit at least one of the selected filters.

Specifications may be seen by pressing `Specifications` buttons.  
Cars can also be compared by selecting desired cars using the checkboxes and then pressing `Compare` button.

`Home` button takes user back to home page.

### Comparison page

Comparison page opens to a new browser tab and shows cars' detailed information in a table so that they are easy to compare.

### Specifications page

Specifications page opens to a new browser tab and shows selected car's detailed information and manufacturer's information.

As a bonus feature user can download the cars information in a text file by pressing `Download text file` button.
