import sys
import os
import whisper

def transcribe_file(input_file, output_file):
    # Load the model
    model = whisper.load_model("base")
    
    try:
        # Transcribe the audio
        result = model.transcribe(input_file)
        
        # Save transcription
        with open(output_file, "w") as f:
            f.write(result["text"])
            
    except Exception as e:
        print(f"Error transcribing {input_file}: {str(e)}")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python3 transcribe.py <input_file> <output_file>")
        sys.exit(1)

    input_file = sys.argv[1]
    output_file = sys.argv[2]

    if not os.path.exists(input_file):
        print(f"Error: Input file {input_file} does not exist")
        sys.exit(1)

    transcribe_file(input_file, output_file) 