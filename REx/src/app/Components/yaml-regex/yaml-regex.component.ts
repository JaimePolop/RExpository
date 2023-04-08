import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import {elementAt, map} from 'rxjs/operators';
import {HttpClient} from '@angular/common/http';
import { parse, stringify } from 'yaml'
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import {MatToolbarModule} from '@angular/material/toolbar'; 
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import {Clipboard} from '@angular/cdk/clipboard';



@Component({
  selector: 'app-yaml-regex',
  templateUrl: './yaml-regex.component.html',
  styleUrls: ['./yaml-regex.component.css']
})
export class YamlRegexComponent {
  
  public all_button:any;
  public misc_button:any;
  public searchFormControl: FormControl<string>;
  public parsedYamlObject: any;
  public paths:any;
  public regular_expresions:any;
  public regular_expresion:any;
  public selected_regular_expresion:any;
  public selected_regular_expresion_regexes:any=[];
  public all_regex: any[] = [];
  public array_all_regex:any= [];
  regex:any;
  displayedColumns = ['Name', 'Regex', 'Case Sensitive', 'False Positives', 'Copy Regex'];
  
  
  constructor(private http: HttpClient, private clipboard: Clipboard) {
    this.searchFormControl = new FormControl();
    this.getYAML().subscribe(data => {

      this.parsedYamlObject = data;
      this.getRegexPath();
      this.getRegex();
    })

  }
  ngOnInit(): void {
    this.getRegex();
    this.searchFormControl = new FormControl();
  }


  public onSearchChanges() {
    this.selected_regular_expresion_regexes = this.selected_regular_expresion["regexes"].filter((x: any) =>
      JSON.stringify(x).toLowerCase().includes(this.searchFormControl.value.toLowerCase())
    );
  }

  public getRegexPath(){
    this.paths = this.parsedYamlObject.paths;
  }

  public getRegex(){
    
    this.regular_expresions = this.parsedYamlObject.regular_expresions;
    //console.log(this.regular_expresions);
    for (let regular_expresion of this.regular_expresions){
      this.all_regex.push(regular_expresion.regexes);
      //console.log(this.all_regex);
      for(let regex of regular_expresion.regexes){
        //console.log(regex);
      }
    }

    //Array for the button All
    this.selected_regular_expresion=[];
    //console.log(this.selected_regular_expresion);
    this.all_regex.forEach((element:any)=>{
      element.forEach((element2:any)=>{
        this.array_all_regex.push(element2);
      })
    });
    this.selected_regular_expresion["regexes"]=this.array_all_regex;
    this.setTableToAll();
  }

  //Parse the YAML
  public getYAML(): Observable<any> {
    return this.http.get("./assets/regex.yaml", {
      observe: 'body',
      responseType: "text"   // This one here tells HttpClient to parse it as text, not as JSON
    }).pipe(
      // Map Yaml to JavaScript Object
      map(yamlString => parse(yamlString))
    );
  }

  //Rest of the buttons
  public setTable(regular_expresion: any){
    this.selected_regular_expresion = regular_expresion;
    this.selected_regular_expresion_regexes = this.selected_regular_expresion["regexes"];
  }
  
  //All button
  public setTableToAll(){
    this.selected_regular_expresion_regexes=this.array_all_regex;

  }

  copyRegex(regular_expresion: any) {
    this.clipboard.copy(regular_expresion);
  }

}
