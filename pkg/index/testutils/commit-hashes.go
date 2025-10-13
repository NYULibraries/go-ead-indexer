package testutils

// ------------------------------------------------------------------------------
// git repo fixture constants shared by cmd/index and pkg/index tests
// ------------------------------------------------------------------------------

/*
	# Commit history from test fixture
	bd5b23c0eeb79ee6603ef2867cc0128afba1f210 Updating README.md
	ce9834cb903f3dbee5c6f6fdd5632e094faf2464 Updating nyuad/ad_mc_019.xml, Deleting file tamwag/tam_143.xml EADID='tam_143', Updating edip/mos_2024.xml, Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Updating akkasah/ad_mc_030.xml
	a436436afbc4e2db3c4adbf124edc5e4c9c6daf7 Updating tamwag/tam_143.xml, Updating cbh/arc_212_plymouth_beecher.xml
	9b1a9a9985119c417719aa4652727d6f49f44c9c Deleting file fales/mss_420.xml EADID='mss_420', Updating fales/mss_420.xml
	57aa05d3e0b2000fe5a6c8d5ac2f7318c1b6da8f Updating fales/mss_420.xml
	62d24fbde8de5186145ba7accc92a52dd5a7f33f Deleting file fales/mss_460.xml EADID='mss_460'
	d33325bcb07672aba7f758d03a65d3c89fb29943 Updating fales/mss_460.xml
	bf88d3b38db71d67914799d27fcafe5bd12a9c11 Deleting file akkasah/ad_mc_030.xml EADID='ad_mc_030', Deleting file cbh/arc_212_plymouth_beecher.xml EADID='arc_212_plymouth_beecher', Deleting file edip/mos_2024.xml EADID='mos_2024', Deleting file fales/mss_420.xml EADID='mss_420', Deleting file fales/mss_460.xml EADID='mss_460', Deleting file nyhs/ms256_harmon_hendricks_goldstone.xml EADID='ms256_harmon_hendricks_goldstone', Deleting file nyhs/ms347_foundling_hospital.xml EADID='ms347_foundling_hospital', Deleting file nyuad/ad_mc_019.xml EADID='ad_mc_019', Deleting file tamwag/tam_143.xml EADID='tam_143'
	ac8ef515cd206a68b6c093cbe894c82a0ff9ca04 Updating akkasah/ad_mc_030.xml, Updating cbh/arc_212_plymouth_beecher.xml, Updating edip/mos_2024.xml, Updating fales/mss_420.xml, Updating fales/mss_460.xml, Updating nyhs/ms256_harmon_hendricks_goldstone.xml, Updating nyhs/ms347_foundling_hospital.xml, Updating nyuad/ad_mc_019.xml, Updating tamwag/tam_143.xml
*/

// hashes from the git-repo fixture (in order of commits)
const AddAllHash = "ac8ef515cd206a68b6c093cbe894c82a0ff9ca04"
const DeleteAllHash = "bf88d3b38db71d67914799d27fcafe5bd12a9c11"
const AddOneHash = "d33325bcb07672aba7f758d03a65d3c89fb29943"
const DeleteOneHash = "62d24fbde8de5186145ba7accc92a52dd5a7f33f"
const DeleteModifyAddHash = "9b1a9a9985119c417719aa4652727d6f49f44c9c"
const AddTwoHash = "a436436afbc4e2db3c4adbf124edc5e4c9c6daf7"
const AddThreeDeleteTwoHash = "ce9834cb903f3dbee5c6f6fdd5632e094faf2464"
const NoEADFilesInCommitHash = "bd5b23c0eeb79ee6603ef2867cc0128afba1f210"

