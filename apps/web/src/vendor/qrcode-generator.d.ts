/**
 * Minimal type declaration for qrcode-generator (MIT, Kazuhiko Arase).
 * The package ships plain CommonJS with no bundled types; only the API the
 * web app uses is declared here (W11 · M-01 QR rendering).
 */
declare module "qrcode-generator" {
  interface QRCode {
    addData(data: string, mode?: string): void;
    make(): void;
    /** Module matrix dimension (rows == cols). */
    getModuleCount(): number;
    /** Whether the module at (row, col) is dark. */
    isDark(row: number, col: number): boolean;
  }

  /** Type number 0 = auto-select the smallest fitting version. */
  type TypeNumber = 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 | 24 | 25 | 26 | 27 | 28 | 29 | 30 | 31 | 32 | 33 | 34 | 35 | 36 | 37 | 38 | 39 | 40;
  type ErrorCorrectionLevel = "L" | "M" | "Q" | "H";

  function qrcode(typeNumber: TypeNumber, errorCorrectionLevel: ErrorCorrectionLevel): QRCode;

  export = qrcode;
}
