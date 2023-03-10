import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { YamlRegexComponent } from './Components/yaml-regex/yaml-regex.component';
import {HttpClientModule} from '@angular/common/http';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import {MatToolbarModule} from '@angular/material/toolbar'; 
import {MatTabsModule} from '@angular/material/tabs';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';
import { YamlRegexDetailsComponent } from './Components/yaml-regex-details/yaml-regex-details.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ReactiveFormsModule } from '@angular/forms';
import { MatCardModule } from "@angular/material/card";
import { MatTableModule } from '@angular/material/table'  

@NgModule({
  declarations: [
    AppComponent,
    YamlRegexComponent,
    YamlRegexDetailsComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    NgbModule,
    MatToolbarModule,
    MatTabsModule,
    NgxDatatableModule,
    BrowserAnimationsModule,
    ReactiveFormsModule,
    MatCardModule,
    MatTableModule
    
    
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
