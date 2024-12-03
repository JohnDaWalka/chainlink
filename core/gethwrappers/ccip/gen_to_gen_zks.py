# utility to convert go_generate.go to go_generate_zks.go

def update_go_generate(file_path):
    updated_lines = []

    with open(file_path, 'r') as file:
        for line in file:
            if "//go:generate" in line:
                parts = line.strip().split()  # Split the line into parts
                bin_path = parts[5]
                # Generate the .zbin path with `.sol` insertion
                zbin_path = bin_path.replace("solc", "zksolc").replace("/v0.8.24/", "/v1.5.6/").replace(".bin", ".sol/") + bin_path.split('/')[-1].replace(".bin", ".zbin")
                parts.append(zbin_path)
                
                parts[3] = "./generation/generate_zks/wrap.go"
                
                line = " ".join(parts) + "\n"
                print(line)
            
            updated_lines.append(line)

    # Write back to the file
    with open(file_path, 'w') as file:
        file.writelines(updated_lines)

# Example usage
update_go_generate('go_generate_zks.go')
