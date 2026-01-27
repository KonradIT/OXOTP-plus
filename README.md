# OXOTP+

<p align="center"><img alt="PICTURE logo" src="img/oxotp-ss2.png" width="450"></p>

<b>FIRST OTP BASED ON M5StickCPlus / 2 for 2FA </b><br><br>
-------

Reimplementation of the original **OXOTP** work of [@IMSHOX](https://github.com/IMSHOX), with added compatibility for *M5StickCPlus* and *M5StickCPlus 2*, and a new web application by [xick](https://github.com/xick).

Forked from [@xick's repo](https://github.com/xick/OXOTP-plus)

Further work done by [@KonradIT](https://github.com/KonradIT) to add:

- ability to input an `otpauth://` URI in the web app: https://github.com/KonradIT/OXOTP-plus/pull/1
- PlatformIO compatibility for building the project
- Captive portal support, so that when a device connects to the AP, it will automatically redirect to the web app: https://github.com/KonradIT/OXOTP-plus/pull/2
- M5 StickC Plus2 LED + tone: https://github.com/KonradIT/OXOTP-plus/pull/3
- Support for M5StickC

<p align="center"><img alt="otps section" src="img/screenshot_otps.png" width="350"></p>

## HOW TO USE

**1.** Clone locally, using PlatformIO, and build the project for your specific board.
**2.**  For the first time use, by going to the wifi settings the **SSID** and **password** of the access point to connect for the configuration will be printed on screen. A popup will appear on your phone to login, similar to captive portal from the airport WiFi. Click on it and the website will be loaded.
**3.** On the **configuration** section of the webapp you can also specify your router ssid and password by changing the wireless mode to host (it will automatically fallback to AP mode if connection fails)
-<p align="center"><img alt="confituration section" src="img/screenshot_config.png" width="200"></p>
**4**. Click '**Sync Time**' to setup the RTC timer of the device (be sure your system / browser time is also correct), this is done automatically by the website config too <br> <br>
**5.** Add OTPs using one of two methods:

**Option A: Paste otpauth:// URI** (Recommended)
- Extract URIs from your existing authenticator app using [**extract_otp_secrets**](https://github.com/scito/extract_otp_secrets)
- Paste the `otpauth://totp/...` URI directly into the parser field
- The form will be auto-populated with the service name, account, and secret

**Option B: Manual Entry**
- Enter the Base32 secret key manually (e.g. `OXXA6YXXVTTP4U25`)
- Fill in the service label and username

<p align="center"><img alt="add section" src="img/screenshot_add.png" width="200"></p>

## HOW TO BUILD (PlatformIO - Recommended):

 - Clone the repo
 - Install [PlatformIO](https://platformio.org/) (VS Code extension or CLI)
 - Open the project folder in VS Code with PlatformIO extension
 - Select the target environment:
   - `m5stick-c` for M5StickC
   - `m5stick-c-plus-2` for M5StickC Plus 2
 - Build: `pio run -e m5stick-c-plus-2`
 - Upload: `pio run -e m5stick-c-plus-2 --target upload`

All dependencies are automatically managed by PlatformIO.

### Modifying the Web Interface:

If you modify `web/index.html`, regenerate the C headers:
```bash
go run tools/generate_headers.go
```


## DEVELOPED FOR DAILY USE:

OXOTP+ is intended to be used as an alternative or backup for your currently 2FA authentication app. Using a dedicated device is great because your phone can be easily lost, stolen, hacked, etc. 
Always keep a backup (even on paper) of your secrets if you intend to use this as your main 2FA system.
<br>
<p align="center"><img alt="PICTURE logo" src="img/oxotp-ss1.png" width="450"></p>
<br>

### MILESTONES
 - ~~Add a battery indicator~~ DONE
 - ~~OTA updates~~  DONE
 - ~~Beautify the UI, add an alternative Light UI~~ DONE
 - ~~Make poweroff timer and screen brightness tweakable~~ DONE
 - ~~Support otpauth:// URI parsing~~ DONE
 - Add API for Exporting secrets 
 - Add **Pincode** for the webapp
 - Embed a **QRCode** reader in the webapp

### COMPATIBILITY
- M5StickC Plus 2
- M5StickC

## License 

This software is licensed under the MIT License, to read the full license <a href="LICENSE" target="_blank">HERE</a>.

## Credits
- Based on https://github.com/IMSHOX
- Pico CSS, a splendid and lightweight css framework (https://github.com/picocss/pico)
- truetype2gfx - https://rop.nl/truetype2gfx/
- and the creators of all the libraries used.

