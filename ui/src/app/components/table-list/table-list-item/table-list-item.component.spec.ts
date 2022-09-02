import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TableListItemComponent } from './table-list-item.component';

describe('TableListItemComponent', () => {
  let component: TableListItemComponent;
  let fixture: ComponentFixture<TableListItemComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TableListItemComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(TableListItemComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
