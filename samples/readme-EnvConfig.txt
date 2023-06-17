{
 "CertificateRootDir" : "/certificates", <-- absolute path, always
 "RootCAdir" : "rootCA", <-- relative path to CertificateRootDir
 "ServerCertsDir" : "servers", <-- relative path to CertificateRootDir
 "CertificatesConfigDir" : "configs", <-- relative path to CertificateRootDir
 "RemoveDuplicates": true <-- this ensures that certificates are unique. should always be set to true
}
