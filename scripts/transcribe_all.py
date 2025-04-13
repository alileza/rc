import os
import sys
import glob
from tqdm import tqdm
import whisper

def transcribe_directory(input_dir, output_dir):
    # Load the model once
    print("Loading Whisper model...")
    model = whisper.load_model("base")
    
    # Get all wav files and sort them
    wav_files = sorted(glob.glob(os.path.join(input_dir, "chunk_*.wav")))
    
    print(f"Found {len(wav_files)} files to transcribe")
    
    # Process each file
    for wav_file in tqdm(wav_files, desc="Transcribing"):
        # Generate output filename
        base_name = os.path.basename(wav_file)
        txt_file = os.path.join(output_dir, base_name.replace(".wav", ".txt"))
        
        try:
            # Transcribe the audio
            result = model.transcribe(wav_file)
            
            # Save transcription
            with open(txt_file, "w") as f:
                f.write(result["text"])
                
        except Exception as e:
            print(f"\nError transcribing {wav_file}: {str(e)}")
            continue

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 transcribe_all.py <input_directory> <output_directory>")
        sys.exit(1)

    input_dir = sys.argv[1]
    output_dir = sys.argv[2]

    if not os.path.exists(input_dir):
        print(f"Error: Input directory {input_dir} does not exist")
        sys.exit(1)

    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    transcribe_directory(input_dir, output_dir) 