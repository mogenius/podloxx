import { Component, EventEmitter, HostListener, OnInit, Output } from '@angular/core';
import { StatsService } from '@lox/services/stats.service';

@Component({
  selector: 'lox-raw-data-modal',
  templateUrl: './raw-data-modal.component.html',
  styleUrls: ['./raw-data-modal.component.scss']
})
export class RawDataModalComponent implements OnInit {
  @HostListener('window:keyup', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (event.key === 'Escape' || event.key === 'f') {
      this.close.emit(true);
    }
  }

  constructor(private readonly _statsService: StatsService) {}

  @Output() close: EventEmitter<boolean> = new EventEmitter<boolean>();

  private _data: string;

  ngOnInit(): void {}

  public toggleRaw(): void {
    this.close.emit(true);
  }

  get data(): string {
    return this._data;
  }
}
