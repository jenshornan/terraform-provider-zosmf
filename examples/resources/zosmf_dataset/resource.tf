resource "zosmf_dataset" "example" {
  name      = "IBMUSER.TATI.TERROF"
  volser    = "zmf046"
  unit      = "3390"
  dsorg     = "PS"
  alcunit   = "TRK"
  primary   = 1
  secondary = 100
  avgblk    = 500
  recfm     = "VB"
  blksize   = 27966
  lrecl     = 800
  content   = <<EOT
test
    EOT
}