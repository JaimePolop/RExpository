import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import {map} from 'rxjs/operators';
import {HttpClient} from '@angular/common/http';
import { parse, stringify } from 'yaml'
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import {MatToolbarModule} from '@angular/material/toolbar'; 
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { YamlRegexDetailsComponent } from '../yaml-regex-details/yaml-regex-details.component';



@Component({
  selector: 'app-yaml-regex',
  templateUrl: './yaml-regex.component.html',
  styleUrls: ['./yaml-regex.component.css']
})
export class YamlRegexComponent {

  
  public searchFormControl: FormControl<string>;
  public parsedYamlObject: any;
  public paths:any;
  public regular_expresions:any;
  public regular_expresion:any;
  public selected_regular_expresion:any;
  public selected_regular_expresion_regexes:any;
  regex:any;
  displayedColumns = ['Name', 'Regex', 'Case Sensitive'];
  
  //tableDetails: YamlRegexDetailsComponent = new YamlRegexDetailsComponent;
  
  constructor(private http: HttpClient) {
    this.searchFormControl = new FormControl();
    this.getYAML().subscribe(data => {

      this.parsedYamlObject = data;
      this.getRegexPath();
      this.getRegex();
      //console.log(data);
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
    //console.log(this.regular_expresions);
  }

  public getRegexPath(){
    this.paths = this.parsedYamlObject.paths;
    //console.log( this.paths);
  }

  public getRegex(){
    this.regular_expresions = this.parsedYamlObject.regular_expresions;
    for (let regular_expresion of this.regular_expresions){
      //console.log(regular_expresion.regexes);
      for(let regex of regular_expresion.regexes){
        //console.log(regex);
      }
    }
    
  }

  public getYAML(): Observable<any> {
    return this.http.get("./assets/regex.yaml", {
      observe: 'body',
      responseType: "text"   // This one here tells HttpClient to parse it as text, not as JSON
    }).pipe(
      // Map Yaml to JavaScript Object
      map(yamlString => parse(yamlString))
    );
  }

  public setTable(regular_expresion: any){
    this.selected_regular_expresion = regular_expresion;
    this.selected_regular_expresion_regexes = this.selected_regular_expresion["regexes"];
    console.log(this.selected_regular_expresion_regexes);
  }


  // private onSearchChanges(){

  // }
}
