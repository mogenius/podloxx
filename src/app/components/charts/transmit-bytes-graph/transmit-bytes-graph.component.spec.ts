import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TransmitBytesGraphComponent } from './transmit-bytes-graph.component';

describe('TransmitBytesGraphComponent', () => {
  let component: TransmitBytesGraphComponent;
  let fixture: ComponentFixture<TransmitBytesGraphComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TransmitBytesGraphComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TransmitBytesGraphComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
