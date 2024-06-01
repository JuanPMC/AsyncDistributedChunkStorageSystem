# Programming Assignment 4

## Video tutorial
You can refer to misc/pocVideo.mp4 for a simple overview of the project.

## Prerequisites

- The deployment scripts will automatically check and install Go.
- A Linux environment is required.

## Usage

1. **Deploy Mesh:**

   Navigate to the 'deployment' directory and run the following command:

   bash ./deploySuperMesh.sh


2. **Deploy tree:**

   bash ./deploySuperTree.sh

   Run the appropriate deployment script for each peer node.
   
  ** Deploy single weak node for testin
  
  bash ./deploySimpleWeaknode.sh

3. **Using the nodes:**
Only use the single weak node or the test scripts to test the deployments, the terminals that deply the entire system have many threads running at the same time ( for testing perpouses ) thay arent ment to be used by a user.


## Troubleshooting

- If everything does not work as expected, try re-deploying everything in order.
- Examine log files for more detailed information on any issues.
