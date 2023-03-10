import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-yaml-regex-details',
  templateUrl: './yaml-regex-details.component.html',
  styleUrls: ['./yaml-regex-details.component.css']
})
export class YamlRegexDetailsComponent {

  @Input() regular_expresions:any;

  regular_expresion:any;

  constructor() {

  }

  ngOnInit(): void {

    for (let regular in this.regular_expresions){
      console.log(regular);
    }
  }
}
