import { ComponentFixture, TestBed } from '@angular/core/testing';

import { BtpComponent } from './btp.component';

describe('BtpComponent', () => {
  let component: BtpComponent;
  let fixture: ComponentFixture<BtpComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [BtpComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(BtpComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
