{
	"Country" : "CA", -> This is the country of origin for the certificate
	"Provice" : "Quebec", -> State of province of origin
	"Locality" : "Blainville", -> City of origin
	"Organization" : "myorg.net", -> Organization of origin
	"OrganizationalUnit" : "myorg", -> Sub-organization of origin
	"CommonName" : "myorg.net root CA", -> The name of the certificate
	"EmailAddresses" : ["certs@myorg.net", "certificates@myorg.net"], -> Array of email addresses responsible for this cert
	"Duration" : 10, -> CA duration, in years
	"KeyUsage" : ["Digital Signature", "Certificate Sign", "CRL Sign"], -> Certificate usage. This here are common values for CAs
	"DNSNames" : ["myorg.net","myorg.com","lan.myorg.net"], -> DNS names assigned to this cert
	"IPAddresses" : ["10.1.1.11", "127.0.0.1"], -> IP addresses assigned to this cert (never a good idea to assign IPs to a CA)
	"CertificateDirectory" : "/tmp/", -> directory where to write the cert
	"CertificateName" : "sample_cert", -> cert filename, no extension to the filename
	"IsCA": true, -> Are we creating a CA or a "normal" server cert ?
	"Comments": ["To see which values to put in the KeyUsage field, see https://pkg.go.dev/crypto/x509#KeyUsage", "Strip off 'KeyUsage' from the const name and there you go.", "", "Please note that this field offers no functionality and is strictly here for documentation purposes"] -> Those won't appear in the certificate file
}