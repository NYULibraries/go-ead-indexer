package testutils

// ------------------------------------------------------------------------------
// git repo fixture constants shared by cmd/index and pkg/index tests
// ------------------------------------------------------------------------------

/*
	# Commit history from test fixture
	6696e0513a6dcb38e14a1da46ac7ba44611c6f90 Updating README.md (HEAD -> master)
	598ce06b5bf534e9dec0db5fd64bee88020c6571 Updating nyuad/ad_mc_019.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml, Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Updating akkasah/ad_mc_030.xml
	50fc07058d893854b2eab1ce6285aa98d6596a16 Updating tamwag/tam_143.xml, Updating cbh/arc_212_plymouth_beecher.xml
	244e53e7827640496ead934516ccb68d5d25cb96 Deleting file fales/mss_420.xml EADID='mss_420', Updating fales/mss_420.xml
	dc63b18f64864f2bdcaffee758e4c590dac8f5ab Updating fales/mss_420.xml
	cb2d1300d7c5572bed7a6f2ec5aa67f023fe087c Deleting file fales/mss_460.xml EADID='mss_460'
	52ac657cc70005670c2ba97c23fba68ce8f1f9de Updating fales/mss_460.xml
	6c82536efc4149599c6d341e34dcc1255131c365 Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143'
	6c814c9836fc2abfa89d49f548fcd9cb11eae78a Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml
*/

// hashes from the git-repo fixture (in order of commits)
const AddAllHash = "6c814c9836fc2abfa89d49f548fcd9cb11eae78a"
const DeleteAllHash = "6c82536efc4149599c6d341e34dcc1255131c365"
const AddOneHash = "52ac657cc70005670c2ba97c23fba68ce8f1f9de"
const DeleteOneHash = "cb2d1300d7c5572bed7a6f2ec5aa67f023fe087c"
const DeleteModifyAddHash = "244e53e7827640496ead934516ccb68d5d25cb96"
const AddTwoHash = "50fc07058d893854b2eab1ce6285aa98d6596a16"
const AddThreeDeleteTwoHash = "598ce06b5bf534e9dec0db5fd64bee88020c6571"
const NoEADFilesInCommitHash = "6696e0513a6dcb38e14a1da46ac7ba44611c6f90"
