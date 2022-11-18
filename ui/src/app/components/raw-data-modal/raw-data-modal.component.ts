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

  @Output() close: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private readonly _statsService: StatsService) {}

  ngOnInit(): void {}

  public toggleRaw(): void {
    this.close.emit(true);
  }

  get data(): any {
    return this._statsService.rawData;
  }
}
