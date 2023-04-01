import { ComponentFixture, TestBed } from '@angular/core/testing';

import { YamlRegexComponent } from './yaml-regex.component';

describe('YamlRegexComponent', () => {
  let component: YamlRegexComponent;
  let fixture: ComponentFixture<YamlRegexComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ YamlRegexComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(YamlRegexComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
