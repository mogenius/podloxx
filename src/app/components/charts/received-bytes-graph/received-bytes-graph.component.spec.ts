import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ReceivedBytesGraphComponent } from './received-bytes-graph.component';

describe('ReceivedBytesGraphComponent', () => {
  let component: ReceivedBytesGraphComponent;
  let fixture: ComponentFixture<ReceivedBytesGraphComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ReceivedBytesGraphComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ReceivedBytesGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
