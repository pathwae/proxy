export default class Certificate {
  issuer = "";
  subject = "";
  notBefore = "";
  notAfter = "";
  commonName = "";
  dnsNames = new Array<string>();
}
