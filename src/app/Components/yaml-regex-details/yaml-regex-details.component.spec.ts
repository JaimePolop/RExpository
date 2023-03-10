import { ComponentFixture, TestBed } from '@angular/core/testing';

import { YamlRegexDetailsComponent } from './yaml-regex-details.component';

describe('YamlRegexDetailsComponent', () => {
  let component: YamlRegexDetailsComponent;
  let fixture: ComponentFixture<YamlRegexDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ YamlRegexDetailsComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(YamlRegexDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
