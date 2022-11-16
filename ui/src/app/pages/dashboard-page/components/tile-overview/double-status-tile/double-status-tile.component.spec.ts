import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DoubleStatusTileComponent } from './double-status-tile.component';

describe('DoubleStatusTileComponent', () => {
  let component: DoubleStatusTileComponent;
  let fixture: ComponentFixture<DoubleStatusTileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ DoubleStatusTileComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(DoubleStatusTileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
