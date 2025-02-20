package main

func handle_command(ochan chan map[string]any, cmd map[string]any) {

	odata := map[string]any{
		"DONE": cmd,
	}
	ochan <- odata

}
