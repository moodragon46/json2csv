# JSON to CSV

These go scripts convert json objects with nested data to csv format and back again. 

## CSV Format

The top row is the "path" of the value in the json file. Each part of the path is of the format type<key>. The type can either be arr (array) or map (object/map).
The bottom row are the values of the json file. Currently two types are supported, string and float.
Here is an example JSON file in csv:

#### JSON

```json
{
    "l1m": {
        "l2a": [
            "l3v",
            7.642
        ],
        "l2v": "level 2 value"
    },
    "l1v": 335.64,
    "morelvl1": "bla bal"
}
```

#### CSV

```
/map<l1m>/map<l2a>/arr<1>,/map<l1m>/map<l2v>,/map<l1m>/map<l2a>/arr<0>,/map<l1v>,/map<morelvl1>
float<7.642>,string<level 2 value>,string<l3v>,float<335.64>,string<bla bal>
```

## Usage

Download the /jsoncsv/ folder of the repository. Main.go is an example file.

`import "./jsoncsv"`

Pass in []byte of the json data to encode:

`csvData, err := jsoncsv.JSON2CSV(jsonBytes)`

...and pass in the []byte of the csv data to decode:

`jsonData, err := jsoncsv.CSV2JSON(csvBytes)`
