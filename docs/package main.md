package main
// import fyne
import (
    "encoding/json"
    "io/ioutil"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)
func main() {
    // student struct, you can use any name
    type Student struct {
        Name  string // name N is capital
        Phone string
    }
    // Now creat a slice/ array to store data
    var myStudentData []Student
    // welcome again
    //lets read data from file and display it in list
    data_from_file, _ := ioutil.ReadFile("myFile_name.txt")
    // file name with extension .txt or .json
    // unmarshall/parse data received from file and save/push in slice
    // 2 argument 1. data source, 2. data slice to store data
    json.Unmarshal(data_from_file, &myStudentData)
    // now we can use this data in our list
    // lets create our list
    // new app
    a := app.New()
    // new title and window
    w := a.NewWindow("CRUD APP")
    // resize window
    w.Resize(fyne.NewSize(400, 400))
    // New label to dispaly name and phone number details
    l_name := widget.NewLabel("...")
    l_name.TextStyle = fyne.TextStyle{Bold: true}
    // for phone number
    l_phone := widget.NewLabel("...")
    // entry widget for name
    e_name := widget.NewEntry()
    e_name.SetPlaceHolder("Enter name here...")
    // entry widget for phone
    e_phone := widget.NewEntry()
    e_phone.SetPlaceHolder("Enter phone here...")
    // submit button
    submit_btn := widget.NewButton("Submit", func() {
        // logic part- store our data in json format
        // let create a struct for our data
        // Get data from entry widget and push to slice
        myData1 := &Student{
            Name:  e_name.Text, // data from name entry widget
            Phone: e_phone.Text,
        }
        // append / push data to our slice
        myStudentData = append(myStudentData, *myData1)
        // * star is very important
        // convert / parse data to json format
        // 3 arguments
        // 1st is our slice data
        // ignore 2nd
        // 3rd is identation, we are using space to indent our data
        final_data, _ := json.MarshalIndent(myStudentData, "", " ")
        // writing data to file
        // it take 3 things
        //file name .txt or .json or anything you want to use
        // data source, final_data in our case
        // writing/reading/execute permission
        ioutil.WriteFile("myFile_name.txt", final_data, 0644)
        /// we are done :)
        // lets test
        // empty input field after data insertion
        e_name.Text = ""
        e_phone.Text = ""
        // refresh name & phone entry object
        e_name.Refresh()
        e_phone.Refresh()
    })
    /// Delete Button
    del_button := widget.NewButton("Del", func() {
        // Create a new Temporary slice
        var TempData []Student // student is the struct we create yesterday
        // now loop through the main slice "myStudentData"
        // and push all the data to TempData slice except the select one
        // here i is the indexs and "e" is the element of slice
        // I don't need index here
        // I will push all the data to tempdata slice
        for _, e := range myStudentData {
            // l_name is the label we create to show details
            // important not equal to is used. Don't append if equal to e.Name
            if l_name.Text != e.Name {
                TempData = append(TempData, e)
            }
        }
        // Now append all the data back to main slice myStudentData
        myStudentData = TempData
        // conver to json and marshall indent
        // 3 element
        // our slice is the data source (1st argument)
        // 2nd is prefix.. we don't need prefix
        // 3rd is the indent. we need a single space
        result, _ := json.MarshalIndent(myStudentData, "", " ")
        // write data to file
        // first argument is file name
        // data source which is result here.
        // last one is permission
        ioutil.WriteFile("myFile_name.txt", result, 0644)
    })
    // list widget
    list := widget.NewList(
        // first argument is item count
        // len() function to get myStudentData slice len
        func() int { return len(myStudentData) },
        // 2nd argument is for widget choice. I want to use label
        func() fyne.CanvasObject { return widget.NewLabel("") },
        // 3rd argument is to update widget with our new data
        func(lii widget.ListItemID, co fyne.CanvasObject) {
            co.(*widget.Label).SetText(myStudentData[lii].Name)
        },
    )
    // update on clicked/ selected
    list.OnSelected = func(id widget.ListItemID) {
        l_name.Text = myStudentData[id].Name
        l_name.Refresh()
        // for phone number
        l_phone.Text = myStudentData[id].Phone
        l_phone.Refresh()
    }
    // paste update button code here
    // Update Button
    update_button := widget.NewButton("Update", func() {
        // Temp slice
        var TempData []Student
        // Data I want to update
        update := &Student{
            Name:  e_name.Text,  // entry name widget
            Phone: e_phone.Text, // entry widget : phone
        }
        // looping through our slice and update the data meeting our criteria
        // _ is to ignore index
        // e is the element/data in the myStudentData slice
        for _, e := range myStudentData {
            // checking data criteria
            if l_name.Text == e.Name {
                // if criteria matched, append updated data else
                TempData = append(TempData, *update)
            } else {
                // else append old data in TempData slice
                TempData = append(TempData, e)
            }
        }
        // first
        myStudentData = TempData
        // convert data to json & write to file
        // first argument is data source // 2nd is prefix // third is indent
        result, _ := json.MarshalIndent(myStudentData, "", " ")
        // write data to file
        // 1st argument is file name
        // 2nd is our data(result) from marshalindent
        // 3rd is the file permission.
        ioutil.WriteFile("myFile_name.txt", result, 0644)
        // refresh & empty entry box and refresh list
        e_name.Text = ""
        e_phone.Text = ""
        e_name.Refresh()
        e_phone.Refresh()
        // refresh list also
        list.Refresh() // if not working cut code and paste after list widget
    })
    // update widget tree and add update button also
    // show and run
    w.SetContent(
        // lets create Hsplite container
        container.NewHSplit(
            // first argument is list of data
            list,
            // 2nd is
            // vbox container
            // show both label
            container.NewVBox(
                l_name, l_phone, e_name,
                e_phone, submit_btn, del_button,
                update_button,
            ),
        ),
    )
    w.ShowAndRun()
}