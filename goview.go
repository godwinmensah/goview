package goview
 
import (
    "encoding/json"
    "fmt"
        "html/template"
    "io/ioutil"
    "net/http"
    "net/url"
 
    "appengine"
        "appengine/urlfetch"
)
 
func init() {
        http.HandleFunc("/", handler)
        http.HandleFunc("/showimage", showimage)
}
 
func handler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, rootForm)
}
 
const rootForm = `
  <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <title>Go View</title>
        <link rel="stylesheet" href="/stylesheets/goview.css">
        <link rel="stylesheet" href="/stylesheets/css/bootstrap.min.css">
      </head>
      <body>
      <header class="container-fluid navigation">
        <h4><a href="#">iSpace Open Data</a></h4>
        <div>
          <ul>
            <li><a href="#">Home</a></li>
            <li><a href="#">Election Data</a></li>
            <li><a href="#">Companies</a></li>
            <li><a href="#">Court Proceedings</a></li>
            <li><a href="#">Sign up</a></li>
            <li><a href="#">Contact</a></li>
          </ul
        </div>
      </header><!--
        <h1><img style="margin-left: 120px;" src="images/gsv.png" alt="iSpace Open Data"</h1>
        <h5>Enter Your Address Location</h5>
        <p>Please enter your address:</p>
        <form style="margin-left: 120px;" action="/showimage" method="post" accept-charset="utf-8">
      <input type="text" name="str" value="Type address..." id="str" />
      <input type="submit" value=".. and see the image!" />
        </form>
       --><div class="container-banner">
            <div class="row"
              <div class="col-lg-1; col-md-1; col-xs-1">
               <img src="images/" alt="">
            </div>
            <h1> There will be slide show of statistics and analytical view of data with texts that will clearly show what this site is all-bout </h1>
          </div>
          <div class="row stacks"
            <div class="col-lg-4; col-md-4; col-xs-4">
              <button type="button" class="btn btn-default">Login</button>
            </div>
          </div>
          <div class="row"
            <div class="col-md-1 col-xs-1">
              <div class="col-md-4; col-xs-4">
                <img src="images/data_warehousing.png" alt="data">
              </div> 
              <div class="col-md-4; col-xs-4">
                <img src="images/companies.jpeg" alt="data">
              </div>
              <div class="col-md-4; col-xs-4 court">
                <img src="images/court.png" alt="data">
              </div>
            </div> 
          </div>
        <div class="footer">
          <div class="row">
            <div class="col-md-4; col-xs-4">
              <p>Ghana Election</p>
                <h5>Presidential Election</h5>
                  <h6>1992 - 2012</h6>
                <h5>Parlimentary Election</h5>
                  <h6>1992 - 2012</h6>
                <h5>District Election</h5>
                  <h6>1992 - 2012</h6>
                <h5>Metropolitan Election</h5>
                  <h6>1992 - 2012</h6>
                <h5>Assemble Election</h5>
                  <h6>1992 - 2012</h6>
            </div>
            <div class="col-md-4; col-xs-4">
              <p>Registered Companies</p>
                <h5>Information Technology</h5>
                <h5>Manifactory Industry</h5>
                <h5>Agricultre Industry</h5>
                <h5>Phamacutucal</h5>
                <h5>Transportation</h5>
                <h5>Building & Construction</h5>
                <h5>Real Estate</h5>
                <h5>Procument</h5>
                <h5>Creative Art</h5>
            </div>
            <div class="col-md-4; col-xs-4">
              <p>Hansards</p>
                <h5>Information Technology</h5>
                <h5>Manifactory Industry</h5>
                <h5>Agricultre Industry</h5>
                <h5>Phamacutucal</h5>
                <h5>Transportation</h5>
                <h5>Building & Construction</h5>
                <h5>Real Estate</h5>
                <h5>Procument</h5>
                <h5>Creative Art</h5>
            </div>
          </div>
        </div>
      </body>
    </html>
`
 
var upperTemplate = template.Must(template.New("showimage").Parse(upperTemplateHTML))
 
func showimage(w http.ResponseWriter, r *http.Request) {
        addr := r.FormValue("str")
 
        safeAddr := url.QueryEscape(addr)
        fullUrl := fmt.Sprintf("http://maps.googleapis.com/maps/api/geocode/json?sensor=false&address=%s", safeAddr)
 
        c := appengine.NewContext(r)
        client := urlfetch.Client(c)
 
        resp, err := client.Get(fullUrl)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }
 
    defer resp.Body.Close()
 
        // Read the content into a byte array
    body, dataReadErr := ioutil.ReadAll(resp.Body)
    if dataReadErr != nil {
        panic(dataReadErr)
    }
 
        res := make(map[string][]map[string]map[string]map[string]interface{}, 0)
 
    json.Unmarshal(body, &res)
 
    lat, _ := res["results"][0]["geometry"]["location"]["lat"]
    lng, _ := res["results"][0]["geometry"]["location"]["lng"]
 
    // %.13f is used to convert float64 to a string
    queryUrl := fmt.Sprintf("http://maps.googleapis.com/maps/api/streetview?sensor=false&size=600x300&location=%.13f,%.13f", lat, lng)
 
        tempErr := upperTemplate.Execute(w, queryUrl)
        if tempErr != nil {
            http.Error(w, tempErr.Error(), http.StatusInternalServerError)
        }
}
 
const upperTemplateHTML = ` 
<!DOCTYPE html>
  <html>
    <head>
      <meta charset="utf-8">
      <title>Display Image</title>
      <link rel="stylesheet" href="/stylesheets/goview.css">              
    </head>
    <body>
      <header class="navigation">
        <ul>
          <li>Home</li>
          <li>About us<li>
          <li>Contact</li>
          </ul>
      </header>
      <h1><img style="margin-left: 120px;" src="images/gsv.png" alt="Street View" />GoView</h1>
      <h2>Image at your Address</h2>
      <img style="margin-left: 120px;" src="{{html .}}" alt="Image" />
    </body>
  </html>
`
