package main

import (
	"fmt"
	"os"
	"testing"
)

func TestHttp(t *testing.T) {
	//err := setupHttpClient()  // load balancing
	err := setupHttpClient2()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to new sdlb: %v\n", err)
		return
	}
	defer closeSDLB()

	TestHttpGetData("click-dark-ch")

	//TestHttpCheckData("click-dark-ch", "aa7e31da-1c6c-11f0-b0d1-8c85907c8cf5", "34,54")

	//TestHttpGetStatusInfo("aa7e31da-1c6c-11f0-b0d1-8c85907c8cf5")

	//TestHttpGetResourceList()

	//TestHttpDelStatusInfo("bba172f4-1c73-11f0-95b9-8c85907c8cf5")

	//TestHttpGetConfig()

	//jsonStr := `{"resources":{"version":"0.0.2000","char":{"type":"","languages":{"chinese":[],"english":[]}},"font":{"type":"load","file_dir":"./gocaptcha/fonts/","file_maps":{"yrdzst_bold":"yrdzst-bold.ttf"}},"shapes_image":{"type":"","file_dir":"","file_maps":null},"master_image":{"type":"load","file_dir":"./gocaptcha/master_images/","file_maps":{"image_01":"image_01.jpg","image_02":"image_02.jpg"}},"thumb_image":{"type":"load","file_dir":"./gocaptcha/thumb_images/","file_maps":{}},"tile_image":{"type":"load","file_dir":"./gocaptcha/tile_images/","file_maps":{"tile_01":"tile_01.png","tile_02":"tile_02.png"},"file_maps_02":{"tile_mask_01":"tile_mask_01.png","tile_mask_02":"tile_mask_02.png"},"file_maps_03":{"tile_shadow_01":"tile_shadow_01.png","tile_shadow_02":"tile_shadow_02.png"}}},"builder":{"click_config_maps":{"click_dark_ch":{"version":"","language":"chinese","master":{"image_size":{"Width":300,"Height":200},"range_length":{"Min":6,"Max":7},"range_angles":[{"Min":20,"Max":35},{"Min":35,"Max":45},{"Min":290,"Max":305},{"Min":305,"Max":325},{"Min":325,"Max":330}],"range_size":{"Min":26,"Max":32},"range_colors":["#fde98e","#60c1ff","#fcb08e","#fb88ff","#b4fed4","#cbfaa9","#78d6f8"],"display_shadow":true,"shadow_color":"#101010","shadow_point":{"X":-1,"Y":-1},"image_alpha":1,"use_shape_original_color":true},"thumb":{"image_size":{"Width":150,"Height":40},"range_verify_length":{"Min":2,"Max":4},"disabled_range_verify_length":false,"range_text_size":{"Min":22,"Max":28},"range_text_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"range_background_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"background_distort":4,"background_distort_alpha":1,"background_circles_num":24,"background_slim_line_num":2,"is_thumb_non_deform_ability":false}},"click_dark_en":{"version":"","language":"english","master":{"image_size":{"Width":300,"Height":200},"range_length":{"Min":6,"Max":7},"range_angles":[{"Min":20,"Max":35},{"Min":35,"Max":45},{"Min":290,"Max":305},{"Min":305,"Max":325},{"Min":325,"Max":330}],"range_size":{"Min":26,"Max":32},"range_colors":["#fde98e","#60c1ff","#fcb08e","#fb88ff","#b4fed4","#cbfaa9","#78d6f8"],"display_shadow":true,"shadow_color":"#101010","shadow_point":{"X":-1,"Y":-1},"image_alpha":1,"use_shape_original_color":true},"thumb":{"image_size":{"Width":150,"Height":40},"range_verify_length":{"Min":2,"Max":4},"disabled_range_verify_length":false,"range_text_size":{"Min":22,"Max":28},"range_text_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"range_background_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"background_distort":4,"background_distort_alpha":1,"background_circles_num":24,"background_slim_line_num":2,"is_thumb_non_deform_ability":false}},"click_default_ch":{"version":"","language":"chinese","master":{"image_size":{"Width":300,"Height":200},"range_length":{"Min":6,"Max":7},"range_angles":[{"Min":20,"Max":35},{"Min":35,"Max":45},{"Min":290,"Max":305},{"Min":305,"Max":325},{"Min":325,"Max":330}],"range_size":{"Min":26,"Max":32},"range_colors":["#fde98e","#60c1ff","#fcb08e","#fb88ff","#b4fed4","#cbfaa9","#78d6f8"],"display_shadow":true,"shadow_color":"#101010","shadow_point":{"X":-1,"Y":-1},"image_alpha":1,"use_shape_original_color":true},"thumb":{"image_size":{"Width":150,"Height":40},"range_verify_length":{"Min":2,"Max":4},"disabled_range_verify_length":false,"range_text_size":{"Min":22,"Max":28},"range_text_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"range_background_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"background_distort":4,"background_distort_alpha":1,"background_circles_num":24,"background_slim_line_num":2,"is_thumb_non_deform_ability":false}},"click_default_en":{"version":"","language":"english","master":{"image_size":{"Width":300,"Height":200},"range_length":{"Min":6,"Max":7},"range_angles":[{"Min":20,"Max":35},{"Min":35,"Max":45},{"Min":290,"Max":305},{"Min":305,"Max":325},{"Min":325,"Max":330}],"range_size":{"Min":34,"Max":48},"range_colors":["#fde98e","#60c1ff","#fcb08e","#fb88ff","#b4fed4","#cbfaa9","#78d6f8"],"display_shadow":true,"shadow_color":"#101010","shadow_point":{"X":-1,"Y":-1},"image_alpha":1,"use_shape_original_color":true},"thumb":{"image_size":{"Width":150,"Height":40},"range_verify_length":{"Min":2,"Max":4},"disabled_range_verify_length":false,"range_text_size":{"Min":34,"Max":48},"range_text_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"range_background_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"background_distort":4,"background_distort_alpha":1,"background_circles_num":24,"background_slim_line_num":2,"is_thumb_non_deform_ability":false}}},"click_shape_config_maps":{"click_shape_default":{"version":"","language":"","master":{"image_size":{"Width":300,"Height":200},"range_length":{"Min":6,"Max":7},"range_angles":[{"Min":20,"Max":35},{"Min":35,"Max":45},{"Min":290,"Max":305},{"Min":305,"Max":325},{"Min":325,"Max":330}],"range_size":{"Min":26,"Max":32},"range_colors":["#fde98e","#60c1ff","#fcb08e","#fb88ff","#b4fed4","#cbfaa9","#78d6f8"],"display_shadow":true,"shadow_color":"#101010","shadow_point":{"X":-1,"Y":-1},"image_alpha":1,"use_shape_original_color":true},"thumb":{"image_size":{"Width":150,"Height":40},"range_verify_length":{"Min":2,"Max":4},"disabled_range_verify_length":false,"range_text_size":{"Min":22,"Max":28},"range_text_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"range_background_colors":["#1f55c4","#780592","#2f6b00","#910000","#864401","#675901","#016e5c"],"background_distort":4,"background_distort_alpha":1,"background_circles_num":24,"background_slim_line_num":2,"is_thumb_non_deform_ability":false}}},"slide_config_maps":{"slide_default":{"version":"","master":{"image_size":{"Width":0,"Height":0},"image_alpha":0},"thumb":{"range_graph_size":{"Min":0,"Max":0},"range_graph_angles":null,"generate_graph_number":0,"enable_graph_vertical_random":false,"range_dead_zone_directions":null}}},"drag_config_maps":{"drag_default":{"version":"","master":{"image_size":{"Width":0,"Height":0},"image_alpha":0},"thumb":{"range_graph_size":{"Min":0,"Max":0},"range_graph_angles":null,"generate_graph_number":0,"enable_graph_vertical_random":false,"range_dead_zone_directions":null}}},"rotate_config_maps":{"rotate_default":{"version":"","master":{"image_square_size":0},"thumb":{"range_angles":null,"range_image_square_sizes":null,"image_alpha":0}}}}}`
	//TestHttpUpdateConfig(jsonStr)

	//rootDir, _ := os.Getwd()
	//filePaths := []string{
	//	path.Join(rootDir, "__example/golang/.cache/file1.pdf"),
	//	path.Join(rootDir, "__example/golang/.cache/file2.jpg"),
	//}
	//files, err := openFiles(filePaths)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Failed to conver file header: %v\n", err)
	//	return
	//}
	//TestHttpUploadResource("test_files", files)
}

func openFiles(filePaths []string) ([]*os.File, error) {
	var files []*os.File

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			closeFiles(files)
			return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
		}
		files = append(files, file)
	}

	return files, nil
}

func closeFiles(files []*os.File) {
	for _, file := range files {
		if file != nil {
			file.Close()
		}
	}
}
